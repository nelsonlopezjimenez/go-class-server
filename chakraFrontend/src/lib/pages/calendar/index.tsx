import { useGetInfoQuery } from '@/api/apiSlice';
import { Heading } from '@chakra-ui/react/typography';

export default function Calendar(): React.JSX.Element {
  const { data, isSuccess } = useGetInfoQuery('calendar');
  return isSuccess ? (
    <>
      <div
        dangerouslySetInnerHTML={{ __html: data }}
        className="markdown-content"
      />
    </>
  ) : (
    <>
      <Heading>There May be a calendar here some day</Heading>
    </>
  );
}
