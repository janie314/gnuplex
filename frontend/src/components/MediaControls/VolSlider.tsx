import ReactSlider from "react-slider";
import "./VolSlider.css";
import { SoundHigh, SoundLow, SoundMin, SoundOff } from "iconoir-react";

function VolSlider(props: {
  vol: number | null;
  setVol: React.Dispatch<React.SetStateAction<number | null>>;
  debouncedVol: number | null;
}) {
  return (
    <div className="slider">
      {(props.vol === 0 || props.vol === null)
        ? <SoundOff />
        : (props.vol <= 33
          ? <SoundMin />
          : (props.vol <= 66 ? <SoundLow /> : <SoundHigh />))}
      <ReactSlider
        className="horizontal-slider"
        thumbClassName="thumb"
        trackClassName="track"
        value={Number(props.vol)}
        max={100}
        onChange={(vol: number) => props.setVol(vol)}
      />
    </div>
  );
}

export { VolSlider };
