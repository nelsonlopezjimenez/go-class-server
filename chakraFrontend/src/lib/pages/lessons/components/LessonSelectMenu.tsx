//@ts-nocheck
import { useGetLessonsQuery } from '@/api/apiSlice';
import {
  AccordionItem,
  AccordionItemContent,
  AccordionItemTrigger,
  AccordionRoot,
} from '@/components/ui/accordion';
import { Button } from '@/components/ui/button';
import {
  DrawerActionTrigger,
  DrawerBackdrop,
  DrawerBody,
  DrawerCloseTrigger,
  DrawerContent,
  DrawerHeader,
  DrawerRoot,
  DrawerTitle,
  DrawerTrigger,
} from '@/components/ui/drawer';
import { Skeleton } from '@/components/ui/skeleton';
import { lessonPrettifier, replaceIds } from '@/helpers/helpers';
import { VStack } from '@chakra-ui/react/stack';
import { useNavigate } from 'react-router';

export const LessonSelectMenu: React.FC = () => {
  const { data, isSuccess } = useGetLessonsQuery('dropdown');
  const navigate = useNavigate();
  return isSuccess ? (
    <>
      <DrawerRoot placement="start">
        <DrawerBackdrop />
        <DrawerTrigger asChild>
          <Button variant="outline" size="sm">
            Select Lesson
          </Button>
        </DrawerTrigger>
        <DrawerContent>
          <DrawerHeader>
            <DrawerTitle>Lesson Select</DrawerTitle>
          </DrawerHeader>
          <DrawerBody>
            <AccordionRoot collapsible defaultValue={['b']}>
              {data.ids.map((id) => (
                <AccordionItem key={id} value={id}>
                  <AccordionItemTrigger>{replaceIds(id)}</AccordionItemTrigger>
                  <AccordionItemContent>
                    <VStack>
                      {data.entities[id].map((lesson) => {
                        return (
                          <DrawerActionTrigger key={id + lesson} asChild>
                            <Button
                              size="sm"
                              width="full"
                              variant="subtle"
                              colorPalette="blue"
                              onClick={() =>
                                navigate(`section/${id}/lesson/${lesson}`)
                              }
                            >
                              {lessonPrettifier(lesson)}
                            </Button>
                          </DrawerActionTrigger>
                        );
                      })}
                    </VStack>
                  </AccordionItemContent>
                </AccordionItem>
              ))}
            </AccordionRoot>
          </DrawerBody>

          <DrawerCloseTrigger />
        </DrawerContent>
      </DrawerRoot>
    </>
  ) : (
    <Skeleton height={12} width={36} />
  );
};
//       <MenuRoot>
//         <MenuTrigger asChild>
//           <Button variant="surface" size="lg">
//             Lesson Select
//           </Button>
//         </MenuTrigger>
//         <MenuContent>
//           {data.ids.map((id) => {
//             return (
//               <MenuRoot
//                 onSelect={(event) => handleClick(event.value)}
//                 key={id}
//                 positioning={{ placement: 'right-start', gutter: 2 }}
//               >
//                 <MenuTriggerItem value="open-recent">
//                   {replaceIds(id)}
//                 </MenuTriggerItem>
//                 {/* <MenuTriggerItem value="open-recent">{id}</MenuTriggerItem> */}
//                 <MenuContent>
//                   {data.entities[id].map((lesson) => {
//                     return (
//                       <MenuItem key={id + lesson} value={`section/${id}/lesson/${lesson}`}>
//                         {lesson}
//                       </MenuItem>
//                     );
//                   })}
//                 </MenuContent>
//               </MenuRoot>
//             );
//           })}
//         </MenuContent>
//       </MenuRoot>
