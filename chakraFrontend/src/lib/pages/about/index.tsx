// import Biography from "./Biography";
// import bioObj from "@/data/bio";
import { useGetInfoQuery } from '@/api/apiSlice';

const About = () => {
  const { data, isSuccess } = useGetInfoQuery('about');

  return isSuccess ? (
    <div
      className="markdown-content"
      dangerouslySetInnerHTML={{ __html: data }}
      style={{ margin: '0 10%' }}
    />
  ) : (
    <>
      <p>It appears you have no information</p>
    </>
  );
};

export default About;

// # EDCC Calendar for 2025/2026

// ## October 2025

// | Sun | Mon | Tue | Wed | Thur | Fri | Sat |
// | --- | --- | --- | --- | ---- | --- | --- |
// |     |     |     | 1   | 2    | 3   | 4   |
// | 5   | 6   | 7   | 8   | 9    | 10  | 11  |
// | 12  | 13  | 14  | 15  | 16   | 17  | 18  |
// | 19  | 20  | 21  | 22  | 23   | 24  | 25  |
// | 26  | 27  | 28  | 29  | 30   | 31  |

// ## November 2025

// | Sun | Mon | Tue | Wed | Thur | Fri | Sat |
// | --- | --- | --- | --- | ---- | --- | --- |
// |     |     |     |     |      |     | 1   |
// | 2   | 3   | 4   | 5   | 6    | 7   | 8   |
// | 9   | 10  | 11  | 12  | 13   | 14  | 15  |
// | 16  | 17  | 18  | 19  | 20   | 21  | 22  |
// | 23  | 24  | 25  | 26  | 27   | 28  | 29  |
// | 30  |

// ## December 2025

// | Sun | Mon | Tue | Wed | Thur | Fri | Sat |
// | --- | --- | --- | --- | ---- | --- | --- |
// |     | 1   | 2   | 3   | 4    | 5   | 6   |
// | 7   | 8   | 9   | 10  | 11   | 12  | 13  |
// | 14  | 15  | 16  | 17  | 18   | 19  | 20  |
// | 21  | 22  | 23  | 24  | 25   | 26  | 27  |
// | 28  | 29  | 30  | 31  |      |     |
