import { APICall } from "../lib/APICall";
import { TimeVolInput } from "./TimeVolInput";

function MediaControls(props: {
  mediadirInputPopup: boolean;
  setMediadirInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
  pos: number;
  setPos: React.Dispatch<React.SetStateAction<number>>;
  vol: number;
  setVol: React.Dispatch<React.SetStateAction<number>>;
  volPosToggle: boolean;
  setVolPosToggle: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  return (
    <div className="flex flex-col justify-center">
      <div className="flex flex-row justify-center mt-1">
        <input
          className="m-1 p-1 border border-solid border-black hover:bg-cyan-300"
          type="button"
          value="⏵"
          onClick={() =>
            APICall.play()
              .then(() => APICall.sleep(2000))
              .then(() => props.setVolPosToggle(!props.volPosToggle))
          }
        />
        <input
          className="m-1 p-1 border border-solid border-black hover:bg-cyan-300"
          type="button"
          value="⏸"
          onClick={() =>
            APICall.pause()
              .then(() => APICall.sleep(2000))
              .then(() => props.setVolPosToggle(!props.volPosToggle))
          }
        />
      </div>
      <div className="flex flex-row justify-center mt-1">
        <input
          type="button"
          className="m-1 p-1 border border-solid border-black hover:bg-cyan-300"
          value="Manage Library"
          onClick={() => {
            props.setMediadirInputPopup(true);
          }}
        />
        <input
          type="button"
          className="m-1 p-1 border border-solid border-black hover:bg-cyan-300"
          value="Cast YouTube"
          onClick={() => {
            const url = window.prompt("YouTube URL:", "") || "";
            APICall.setOriginMedia(url);
          }}
        />
      </div>
      <div className="flex flex-row justify-center mt-1">
        <TimeVolInput
          rawtime={props.pos}
          setRawtime={props.setPos}
          type="time"
        />
      </div>
      <div className="flex flex-row justify-center items-center mt-1">
        <span className="mx-1">Vol</span>
        <input
          type="range"
          min={0}
          max="100"
          value={props.vol}
          className="range range-xs"
          // @ts-ignore
          onChange={(e) => props.setVol(e.target.value)}
          onMouseUp={() => APICall.setOriginVol(props.vol)}
        />
        <span className="mx-1">{props.vol}</span>
      </div>
    </div>
  );
}

export { MediaControls };
