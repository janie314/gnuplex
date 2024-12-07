import pause from "../assets/pause.svg";
import play from "../assets/play.svg";
import { API, type SubTrack } from "../lib/API";
import { SubSelector } from "./SubSelector";

function MediaControls(props: {
  mediadirInputPopup: boolean;
  setMediadirInputPopup: React.Dispatch<React.SetStateAction<boolean>>;
  setCastPopup: React.Dispatch<React.SetStateAction<boolean>>;
  startPos: number;
  pos: number;
  setPos: React.Dispatch<React.SetStateAction<number>>;
  timeRemaining: number;
  vol: number;
  setVol: React.Dispatch<React.SetStateAction<number>>;
  subs: SubTrack[] | null;
}) {
  console.log(play);
  return (
    <div className="flex flex-row flex-wrap items-center justify-center content-baseline p-1">
      <div className="mr-1">
        <button
          type="button"
          className="p-2 w-8 border border-solid border-black hover:bg-cyan-200"
          onClick={() => API.play()}
        >
          <img src={play} alt="Play icon" />
        </button>
      </div>
      <div className="mr-2">
        <button
          type="button"
          className="p-2 w-8 border border-solid border-black hover:bg-cyan-200 hover:bg-cyan-200 hover:dark:bg-cyan-950 dark:bg-stone-950 dark:text-slate-300"
          onClick={() => API.pause()}
        >
          <img src={pause} alt="Pause icon" />
        </button>
      </div>
      <div className="flex flex-col max-w-sm grow p-1">
        <div className="flex flex-row items-center">
          <span className="mx-1 dark:text-slate-300">Pos</span>
          <input
            type="range"
            min={0}
            max={props.startPos + props.timeRemaining}
            value={props.pos}
            className="range range-xs dark:range-info"
            onChange={(e) => props.setPos(e.target.valueAsNumber)}
            onMouseUp={() => API.setPos(props.pos)}
            onTouchCancel={() => API.setPos(props.pos)}
          />
          <span className="mx-1">{timeFormat(props.pos)}</span>
        </div>

        <div className="flex flex-row mt-3">
          <span className="mx-1 dark:text-slate-300">Vol</span>
          <input
            type="range"
            min={0}
            max={120}
            value={props.vol}
            className="range range-xs"
            onChange={(e) => props.setVol(e.target.valueAsNumber)}
            onMouseUp={() => API.setVol(props.vol)}
            onTouchCancel={() => API.setVol(props.vol)}
          />
          <span className="mx-1 dark:text-slate-300">{props.vol}</span>
        </div>
      </div>
      <div className="flex flex-row justify-center mt-3 p-1">
        <SubSelector subs={props.subs} />
        <input
          type="button"
          className="mr-1 p-1 border border-solid border-black hover:bg-cyan-200 hover:dark:bg-cyan-950 dark:bg-stone-950 dark:text-slate-300"
          value="Manage Library"
          onClick={() => {
            props.setMediadirInputPopup(true);
          }}
        />
        <input
          type="button"
          className="p-1 border border-solid border-black hover:bg-cyan-200 hover:dark:bg-cyan-950 dark:bg-stone-950 dark:text-slate-300"
          value="Cast URL"
          onClick={() => props.setCastPopup(true)}
        />
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
