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
    <div className="flex flex-col">
      <div className="controlgroup">
        <input
          className="play-button"
          type="button"
          value="⏵"
          onClick={() =>
            APICall.play()
              .then(() => APICall.sleep(2000))
              .then(() => props.setVolPosToggle(!props.volPosToggle))
          }
        />
        <input
          className="pause-button"
          type="button"
          value="⏸"
          onClick={() =>
            APICall.pause()
              .then(() => APICall.sleep(2000))
              .then(() => props.setVolPosToggle(!props.volPosToggle))
          }
        />
      </div>
      <div className="controlgroup">
        <input
          type="button"
          value="Manage Library"
          onClick={() => {
            props.setMediadirInputPopup(true);
          }}
        />
        <input
          type="button"
          value="Cast YouTube"
          onClick={() => {
            const url = window.prompt("YouTube URL:", "") || "";
            APICall.setOriginMedia(url);
          }}
        />
      </div>
      <div className="controlgroup">
        <TimeVolInput
          rawtime={props.pos}
          setRawtime={props.setPos}
          type="time"
        />
      </div>
      <div className="controlgroup">
        <TimeVolInput vol={props.vol} setVol={props.setVol} type="vol" />
      </div>
    </div>
  );
}

export { MediaControls };
