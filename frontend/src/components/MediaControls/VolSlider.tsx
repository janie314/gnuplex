import ReactSlider from "react-slider";
import "./VolSlider.css";
import { SoundHigh, SoundLow, SoundMin, SoundOff } from "iconoir-react";
import { useEffect, useState } from "react";
import { useDebounce } from "usehooks-ts";
import { APICall } from "../../lib/API.ts";

function VolSlider() {
  const [flush, setFlush] = useState(false);
  const [trueVol, setTrueVol] = useState<number | null>(null);
  const [vol, setVol] = useState<number | null>(null);
  const debouncedVol = useDebounce(vol, 500);

  useEffect(() => {
    APICall.vol().then((vol) => {
      if (vol !== null) {
        setVol(vol);
      }
    });
  }, [flush]);

  useEffect(() => {
    if (vol !== null) {
      APICall.setVol(debouncedVol as number).then(() => setFlush(!flush));
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
        value={trueVol !== null ? trueVol : (vol || 0)}
        max={100}
        onChange={(vol: number) => {
          setVol(vol);
          setTrueVol(null);
        }}
      />
    </div>
  );
}

export { VolSlider };
