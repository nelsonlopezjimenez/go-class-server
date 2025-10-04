type IframeProps = {
  src: string;
};

const Iframe = (props: IframeProps): React.JSX.Element => {
  return (
    <iframe
      title={props.src}
      style={{ width: '100%', height: '60vh', border: 'none' }}
      src={`http://localhost:22022/lessons/${props.src}`}
    />
  );
};

export default Iframe;
