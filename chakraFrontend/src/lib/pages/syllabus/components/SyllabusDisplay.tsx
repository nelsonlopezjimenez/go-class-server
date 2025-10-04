import { useGetInfoQuery } from '@/api/apiSlice';
import { useEffect } from 'react';
import { useParams } from 'react-router';

const SyllabusDisplay = () => {
  //   const [courseName, _setCourseName] = useState('index');
  const params = useParams();
  const {
    data: info,
    isSuccess: infoSuccess,
    isFetching: _infoFetching,
  } = useGetInfoQuery(`syllabi/${params.section || 'index'}`);
  useEffect(() => {
    console.log(params);
  }, [params]);
  return (
    <>
      {infoSuccess && (
        <div
          className="markdown-content"
          id="content-display"
          dangerouslySetInnerHTML={{ __html: info }}
        />
      )}
    </>
  );
};
export default SyllabusDisplay;
