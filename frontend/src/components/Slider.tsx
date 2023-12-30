import ReactSlider from "react-slider";
import "./Slider.css";

function Slider(props: {
  vol: number | null;
  setVol: React.Dispatch<React.SetStateAction<number | null>>;
  debouncedVol: number | null;
}) {
  return (
    <>
      <span>{props.debouncedVol}</span>
      <ReactSlider
        className="horizontal-slider"
        thumbClassName="thumb"
        trackClassName="track"
        value={Number(props.vol)}
        onChange={(vol: number) => props.setVol(vol)}
      />
    </>
  );
}

export { Slider };
