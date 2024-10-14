import { APICall } from "../lib/APICall";

enum TimeUnit {
  secs = 0,
  mins = 1,
  hrs = 2,
}

function TimeVolInput(props: {
  rawtime?: number;
  setRawtime?: React.Dispatch<React.SetStateAction<number>>;
  type: string;
}) {
  if (
    props.type === "time" &&
    props.rawtime !== undefined &&
    props.setRawtime !== undefined
  ) {
    const secs = props.rawtime % 60;
    const mins = ((props.rawtime - secs) % 3600) / 60;
    const hrs = (props.rawtime - 60 * mins - secs) / 3600;

    function setTime(val: string, unit: TimeUnit) {
      const n = Number(val) || 0;
      if (props.rawtime === undefined || props.setRawtime === undefined) {
        return;
      }
      if (unit === TimeUnit.hrs) {
        if (n < 0) {
          props.setRawtime(60 * mins + secs);
        } else {
          props.setRawtime(3600 * n + 60 * mins + secs);
        }
      } else if (unit === TimeUnit.mins) {
        if (n < 0 || n > 59) {
          props.setRawtime(3600 * hrs + secs);
        } else {
          props.setRawtime(3600 * hrs + 60 * n + secs);
        }
      } else {
        if (n < 0 || n > 59) {
          props.setRawtime(3600 * hrs + 60 * mins);
        } else {
          props.setRawtime(3600 * hrs + 60 * mins + n);
        }
      }
    }

    return (
      <div className="flex flex-row">
        <span className="flex items-center mr-2">Pos</span>
        <input
          className="w-[3ch] border border-black p-1"
          type="text"
          value={Number(hrs).toString().padStart(2, "0")}
          onChange={(e) => setTime(e.target.value, TimeUnit.hrs)}
        />
        <span className="flex items-center mx-1">:</span>
        <input
          className="w-[3ch] border border-black p-1"
          type="text"
          value={Number(mins).toString().padStart(2, "0")}
          onChange={(e) => setTime(e.target.value, TimeUnit.mins)}
        />
        <span className="flex items-center mx-1">:</span>
        <input
          className="w-[3ch] border border-black p-1 mr-3"
          type="text"
          value={Number(secs).toString().padStart(2, "0")}
          onChange={(e) => setTime(e.target.value, TimeUnit.secs)}
        />
        <input
          type="button"
          className="m-1 p-1 border border-solid border-black hover:bg-cyan-300"
          value="Set"
          min={0}
          max={250}
          onClick={(e) => {
            APICall.setOriginPos(props.rawtime as number);
          }}
        />
      </div>
    );
  }
  return <></>;
}

export { TimeVolInput };
