import { index, route } from '@react-router/dev/routes';

export default [
  index('lib/pages/home/index.tsx'),
  route('offline-links', 'lib/pages/offlineLinks/index.tsx'),
  route('course-content', 'lib/pages/lessons/index.tsx', [
    route(
      'section/:section/lesson/:lesson',
      'lib/pages/lessons/components/LessonDisplay.tsx',
    ),
  ]),
  route('section', 'lib/pages/syllabus/index.tsx', [
    route(':section?', 'lib/pages/syllabus/components/SyllabusDisplay.tsx'),
  ]),
  route('about', 'lib/pages/about/index.tsx'),
  route('calendar', 'lib/pages/calendar/index.tsx'),
  route('lan-links', 'lib/pages/lan-links/index.tsx'),
  route('*', 'lib/pages/404.tsx'),
];
