import ReactSlider from "react-slider";
import "./VolSlider.css";
import { SoundHigh, SoundLow, SoundMin, SoundOff } from "iconoir-react";
import { useEffect, useState } from "react";
import { useDebounce } from "usehooks-ts";
import { APICall } from "../../lib/API.ts";

function VolSlider(props: {
  trueVol: number | null;
  setTrueVol: React.Dispatch<React.SetStateAction<number | null>>;
}) {
  const [vol, setVol] = useState<number | null>(null);
  const debouncedVol = useDebounce(vol, 500);

  useEffect(() => {
    if (vol !== null) {
      APICall.setVol(debouncedVol as number);
    }
  }, [debouncedVol]);

  return (
    <div className="slider">
      {(vol === 0 || vol === null)
        ? <SoundOff />
        : (vol <= 33
          ? <SoundMin />
          : (vol <= 66 ? <SoundLow /> : <SoundHigh />))}
      <ReactSlider
        className="horizontal-slider"
        thumbClassName="thumb"
        trackClassName="track"
        value={props.trueVol !== null ? props.trueVol : (vol || 0)}
        max={100}
        onChange={(vol: number) => {
          setVol(vol);
          props.setTrueVol(null);
        }}
      />
    </div>
  );
}

export { VolSlider };
