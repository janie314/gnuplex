import { useEffect, useState } from "react";
import { APICall } from "../lib/APICall";
import "../App.css";

function Medialist(props: {
  medialist: string[];
  subtitle: string;
  setMedia: React.Dispatch<React.SetStateAction<string>>;
}) {
  return (
    <div className="moviegroup">
      <span className="subtitle">{props.subtitle}</span>
      {props.medialist.map((mediafile: string, i: number) => (
        <input
          type="button"
          className="mediafile"
          key={mediafile}
          onClick={(e) => {
            APICall.setOriginMedia(mediafile)
              .then(() => APICall.sleep(2000))
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
