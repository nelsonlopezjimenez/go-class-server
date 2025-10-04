import { useDownloadWebsiteMutation } from '@/api/apiSlice';
import {
  DialogActionTrigger,
  DialogBody,
  DialogCloseTrigger,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Button } from '@chakra-ui/react/button';
import { Em, Heading, Strong, Text } from '@chakra-ui/react/typography';
import type React from 'react';
import { MdCancel, MdCheck } from 'react-icons/md';

type props = {
  domain: string;
};

const WebsiteDownloadDialog: React.FC<props> = ({ domain }) => {
  const [download, { isSuccess, isLoading, isError }] =
    useDownloadWebsiteMutation();
  const handleDownload = (dom) => {
    download(dom);
  };

  return (
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Are you sure?</DialogTitle>
      </DialogHeader>
      <DialogBody>
        {!isLoading && !isSuccess && (
          <>
            <Text>Clicking the install button will download {domain}</Text>
            <Strong>Be prepared to wait up to several minutes</Strong>
          </>
        )}
        {isLoading && (
          <>
            <Em>Please be patient. This may take several minutes. </Em>
            <Text>
              You must remain connected to the docking station while the site is
              being downloaded
            </Text>
          </>
        )}
        {isSuccess && (
          <>
            <Heading>SUCCESS</Heading>
            <Text>Click the green check to refresh your websites list</Text>
          </>
        )}
      </DialogBody>
      <DialogFooter>
        {isSuccess ? (
          <DialogActionTrigger asChild>
            <Button colorPalette="green">
              <MdCheck />
            </Button>
          </DialogActionTrigger>
        ) : isError ? (
          <>
            <DialogActionTrigger asChild>
              <Button variant="outline">Cancel</Button>
            </DialogActionTrigger>
            <Button
              loading={isLoading}
              colorPalette="red"
              onClick={() => handleDownload(domain)}
            >
              <MdCancel />
            </Button>
          </>
        ) : (
          <>
            <DialogActionTrigger asChild>
              <Button variant="outline">Cancel</Button>
            </DialogActionTrigger>
            <Button
              loading={isLoading}
              colorPalette="blue"
              onClick={() => handleDownload(domain)}
            >
              Install
            </Button>
          </>
        )}
      </DialogFooter>
      <DialogCloseTrigger />
    </DialogContent>
  );
};

export default WebsiteDownloadDialog;
