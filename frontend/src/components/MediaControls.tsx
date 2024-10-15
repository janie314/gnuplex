import { APICall } from "../lib/APICall";

function MediaControls(props: {
  mediadirInputPopup: boolean;
  setMediadirInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
  startPos: number;
  pos: number;
  timeRemaining: number;
  setPos: React.Dispatch<React.SetStateAction<number>>;
  vol: number;
  setVol: React.Dispatch<React.SetStateAction<number>>;
  volPosToggle: boolean;
  setVolPosToggle: React.Dispatch<React.SetStateAction<boolean>>;
}) {
  return (
    <div className="flex flex-col justify-center p-1">
      <div className="flex flex-row justify-center mt-3">
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
      <div className="flex flex-row justify-center mt-3">
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
      <div className="flex flex-row justify-center items-center mt-3">
        <span className="mx-1">Pos</span>
        <input
          type="range"
          min={0}
          max={props.startPos + props.timeRemaining}
          value={props.pos}
          className="range range-xs"
          // @ts-ignore
          onChange={(e) => props.setPos(e.target.value)}
          onMouseUp={() => APICall.setOriginPos(props.pos)}
        />
        <span className="mx-1">{timeFormat(props.pos)}</span>
      </div>

      <div className="flex flex-row justify-center items-center mt-3">
        <span className="mx-1">Vol</span>
        <input
          type="range"
          min={0}
          max={120}
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

function timeFormat(n: number) {
  const secs = Math.floor(n % 60);
  const mins = Math.floor(((n - secs) % 3600) / 60);
  const hrs = Math.floor((n - 60 * mins - secs) / 3600);
  if (hrs === 0) {
    return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
  }
  return `${hrs.toString().padStart(2, "0")}:${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
}

export { MediaControls };
