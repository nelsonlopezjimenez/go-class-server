import { lessonSorter } from '@/helpers/helpers';
import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import markdownit from 'markdown-it';

export const md = markdownit({ html: true, typographer: true });
interface normalizedData {
  ids: Array<string>;
  entities: any;
}
export const apiSlice = createApi({
  reducerPath: 'apiSlice',
  baseQuery: fetchBaseQuery({ baseUrl: '/' }),
  tagTypes: ['siteList'],
  endpoints: (builder) => ({
    getLinks: builder.query({
      query: () => 'api/links',
      providesTags: ['siteList'],
    }),
    getInfo: builder.query({
      query: (param) => `api/information/${param}`,
      transformResponse(res) {
        return md.render(res);
      },
    }),
    // getMd: builder.query({ query: () => 'api/markdown' }),
    getLessons: builder.query({
      query: () => 'api/lessons',
      transformResponse(res) {
        /**
         * @example transformed response before string replacement
         * {"ids":["q1","q2","q3","q4"],
         * "entities":{
         * "q1":["index.md","lesson1.md","lesson2.md","lesson3.md"],
         * "q2":["index.md","lesson1.md","lesson2.md","lesson3.md"],
         * "q3":["index.md","lesson1.md","lesson2.md","lesson3.md"],
         * "q4":["index.md","lesson1.md","lesson2.md","lesson3.md"]}}
         *
         * @return normalized data
         */
        const normalizedData: normalizedData = { ids: [], entities: {} };
        for (const k in res) {
          normalizedData.entities[k] = res[k].map((doc: string) => {
            let tmp = doc;
            tmp = tmp.replace('.md', '');
            return tmp;
          });
          //This sorts the entities array using a custom sorting function that uses any integer at the end of the string to sort values
          //otherwise it sorts alphabetically i.e: a, b, week1, week 2, week10
          normalizedData.entities[k] = lessonSorter(normalizedData.entities[k]);
        }
        normalizedData.ids.push(...Object.keys(res));
        return normalizedData;
      },
    }),
    getLesson: builder.query<string, any>({
      query: (lesson) => `/raw/lessons/${lesson}`,
      transformResponse: (response: string) => {
        return md.render(response);
        // return response
      },
    }),
    downloadWebsite: builder.mutation<string, any>({
      query: (domain) => ({
        url: `/api/git/update/${domain}`,
        method: 'GET',
      }),
      invalidatesTags: ['siteList'],
    }),
  }),
});

export const {
  useGetLessonQuery,
  useGetLinksQuery,
  useGetInfoQuery,
  // useGetMdQuery,
  useGetLessonsQuery,
  useDownloadWebsiteMutation,
} = apiSlice;
