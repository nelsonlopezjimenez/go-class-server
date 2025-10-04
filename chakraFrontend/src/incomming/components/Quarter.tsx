import { useParams } from 'react-router';
import { QUARTER_INFO } from '../data/courseInfo';
import type { QuarterInfo } from '../data/courseInfo';
import CourseDetails from './CourseDetails';

type props = {
  placeholder: string;
};

const Quarter: React.FC<props> = (_props) => {
  const { quarter } = useParams();
  let dataset: Array<QuarterInfo>;

  const firstToUpper = (str: string): string => {
    const firstUpper = str[0].toUpperCase();
    return str.replace(str[0], firstUpper);
  };

  switch (quarter) {
    case '1':
      dataset = QUARTER_INFO.filter((course) => course.quarter === 'one');
      // "one";
      break;
    case '2':
      dataset = QUARTER_INFO.filter((course) => course.quarter === 'two');
      break;
    case '3':
      dataset = QUARTER_INFO.filter((course) => course.quarter === 'three');
      break;
    case '4':
      dataset = QUARTER_INFO.filter((course) => course.quarter === 'four');
      break;
    default:
      return (
        <>
          <h1>The url you followed points to an invalid course.</h1>
        </>
      );
  }

  const courseBoxes = dataset.map((course) => (
    <CourseDetails key={JSON.stringify(course)} details={course} />
  ));

  return (
    <>
      <h1>Quarter {firstToUpper(dataset[0].quarter)}</h1>
      {courseBoxes}
    </>
  );
};

export default Quarter;
