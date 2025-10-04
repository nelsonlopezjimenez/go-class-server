import { useGetLessonQuery } from '@/api/apiSlice';
import { useParams } from 'react-router';

const LessonDisplay = () => {
  // const [lessonName, _setLessonName] = useState('index');
  const params = useParams();
  const {
    data: lesson,
    isSuccess: lessonSuccess,
    isFetching: _lessonFetching,
  } = useGetLessonQuery(`${params.section}/${params.lesson}` || 'index');
  return (
    <>
      {/* {JSON.stringify(params)}
      {`/raw/lessons/${params.section}/${params.lesson}`} */}
      {lessonSuccess && (
        <div
          className="markdown-content"
          id="content-display"
          dangerouslySetInnerHTML={{ __html: lesson }}
        />
      )}
    </>
  );
};
export default LessonDisplay;
