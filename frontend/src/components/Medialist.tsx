import { useEffect, useState } from "react";
import { APICall } from "../lib/API";
import "../App.css";

function Medialist(props: { medialist: string[]; subtitle: string }) {
  return (
    <div className="moviegroup">
      <span className="subtitle">{props.subtitle}</span>
      {props.medialist.map((mediafile: string, i: number) => (
        <a
          className="mediafile"
          key={i}
          href="#"
          onClick={() => APICall.setOriginMedia(mediafile)}
        >
          {mediafile}
        </a>
      ))}
    </div>
  );
}

export { Medialist };
