import { Outlet } from 'react-router';
import { LessonSelectMenu } from './components/LessonSelectMenu';

// import '@/lib/styles/md.css'

function Lesson(): React.JSX.Element {
  return (
    <>
      <LessonSelectMenu />
      <Outlet />
    </>
  );
}

export default Lesson;
