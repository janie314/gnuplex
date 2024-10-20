import { API } from "../lib/API";

function Medialist(props: {
  medialist: string[];
  subtitle: string;
  setMedia: React.Dispatch<React.SetStateAction<string>>;
}) {
  return (
    <div className="flex flex-col mb-2 pl-2 whitespace-pre=wrap">
      <h1 className="m-y-2/100 text-lg font-bold">{props.subtitle}</h1>
      {props.medialist.map((mediafile: string, i: number) => (
        <input
          type="button"
          className="text-left hover:bg-cyan-200 p-1 border-lightgray border-t-2 whitespace-normal break-words"
          key={mediafile}
          onClick={(e) => {
            API.setOriginMedia(mediafile)
              .then(() => API.sleep(2000))
              .then(() => {
                props.setMedia(mediafile);
              });
          }}
          value={mediafile}
        />
      ))}
    </div>
  );
}

export { Medialist };
