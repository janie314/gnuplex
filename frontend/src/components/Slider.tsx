import ReactSlider from "react-slider";
import "./Slider.css";

function Slider() {
  return (
    <ReactSlider
      className="horizontal-slider"
      thumbClassName="thumb"
      trackClassName="track"
    />
  );
}

export { Slider };
