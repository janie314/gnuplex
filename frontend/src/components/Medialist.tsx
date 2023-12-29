import { SyntheticEvent } from "react";
import { APICall } from "../lib/API.ts";
import "../App.css";

function Medialist(
  props: {
    medialist: string[];
    subtitle: string;
    setMedia: React.Dispatch<React.SetStateAction<string>>;
  },
) {
  return (
    <div className="moviegroup">
      <span className="subtitle">{props.subtitle}</span>
      {props.medialist.map((mediafile: string, i: number) => (
        <a
          className="mediafile"
          key={i}
          href="#"
          onClick={(e: SyntheticEvent) => {
            e.preventDefault();
            e.stopPropagation();
            APICall.setOriginMedia(mediafile).then(() => APICall.sleep(2000))
              .then(() => {
                props.setMedia(mediafile);
              });
          }}
        >
          {mediafile}
        </a>
      ))}
    </div>
  );
}

export { Medialist };
