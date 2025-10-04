import _ from 'lodash';
import type { Location } from 'react-router';

export const replaceIds = (id: string) => {
  const idsToReplace = {
    q1: 'Quarter 1',
    q2: 'Quarter 2',
    q3: 'Quarter 3',
    q4: 'Quarter 4',
  };
  //@ts-ignore
  return Object.keys(idsToReplace).includes(id) ? idsToReplace[id] : id;
};

export const locationChecker = (
  location: Location,
  setYouAreHere: (a: string) => void,
) => {
  const locationPathname =
    location.pathname === '/'
      ? 'Home'
      : _.startCase(location.pathname.split('/')[1]);
  setYouAreHere(locationPathname);
};

const matchADigit = /\d/;
const sortFunction = (a, b) => {
  console.log(
    Number.parseInt(
      a.search(matchADigit) !== -1 ? a.slice(a.search(matchADigit)) : '0',
    ),
  );
  console.log(
    Number.parseInt(
      b.search(matchADigit) !== -1 ? b.slice(b.search(matchADigit)) : '0',
    ),
  );
  const cleanB = Number.parseInt(
    b.search(matchADigit) !== -1 ? b.slice(b.search(matchADigit)) : '0',
  );
  const cleanA = Number.parseInt(
    a.search(matchADigit) !== -1 ? a.slice(a.search(matchADigit)) : '0',
  );
  //   return a-b
  return cleanA - cleanB;
};
export const lessonSorter = (lessons: Array<string>) => {
  return lessons.sort().sort((a, b) => sortFunction(a, b));
};

export const lessonPrettifier = (str) => {
  return _.startCase(str);
};
