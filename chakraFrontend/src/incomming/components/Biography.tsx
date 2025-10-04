export type Bio = {
  name: string;
  role: string;
  interests: Array<string>;
  bioTxt: string;
};

const Biography = (props: Bio): React.JSX.Element => {
  return (
    <div style={{ width: '80%', margin: '0 10%' }}>
      <h3 style={{ textDecoration: 'underline', textAlign: 'center' }}>
        {props.name}
      </h3>
      <h4>{props.role}</h4>
      <p>Interests: {props.interests.join(', ')}.</p>
      <p>{props.bioTxt}</p>
    </div>
  );
};

export default Biography;
