import { useEffect, useState } from "react";
import "./Popup.css";
import { API } from "../../lib/API";

const filters = [
  { label: "Black and White", value: "bw" },
  { label: "Grainy", value: "grainy" },
  { label: "8-Bit", value: "8bit" },
  { label: "Mirror", value: "mirror" },
  { label: "Sepia", value: "sepia" },
  { label: "Psychedelic", value: "psychedelic" },
  { label: "Tron", value: "tron" },
  { label: "#NoFilter", value: "neutral" },
];

function SettingsPopup(props: { visible: boolean; closeHook: () => void }) {
  const [subDelay, setSubDelay] = useState(0);
  const [subDelayText, setSubDelayText] = useState("0");
  const [isExternalSub, setIsExternalSub] = useState(false);

  const refreshSubDelay = async () => {
    try {
      const delay = await API.getSubDelay();
      setSubDelay(delay);
      setSubDelayText(delay.toString());
    } catch {
      setSubDelay(0);
      setSubDelayText("0");
    }
  };

  useEffect(() => {
    refreshSubDelay();
    API.getSubTracks()
      .then((tracks) => {
        if (tracks) {
          const selected = tracks.find((t) => t.selected);
          setIsExternalSub(selected ? selected.external : false);
        } else {
          setIsExternalSub(false);
        }
      })
      .catch(() => setIsExternalSub(false));
  }, [props.visible]); // eslint-disable-line react-hooks/exhaustive-deps

  if (props.visible) {
    return (
      <div className="popup bg-white dark:bg-stone-800 m-5 min-w-80 p-12">
        <div className="flex flex-col gap-6">
          <select
            className="btn-standard w-auto"
            onChange={(e) => API.setFilter(e.target.value)}
            defaultValue="neutral"
          >
            <option value="" disabled>
              Select Filter
            </option>
            {filters.map((f) => (
              <option key={f.value} value={f.value}>
                {f.label}
              </option>
            ))}
          </select>
          <label className="flex items-center gap-2 text-black dark:text-white text-sm">
            Sub Delay (s):
            <input
              type="text"
              inputMode="numeric"
              pattern="-?[0-9]*\.?[0-9]*"
              value={subDelayText}
              onChange={(e) => {
                const val = e.target.value;
                if (/^-?[0-9]*\.?[0-9]*$/.test(val)) {
                  setSubDelayText(val === "" || val === "-" ? val : val);
                }
              }}
              onBlur={(e) => {
                const parsed = parseFloat(e.target.value);
                const value = isNaN(parsed) ? 0 : parsed;
                setSubDelay(value);
                setSubDelayText(value.toString());
                API.setSubDelay(value);
              }}
              className="btn-standard w-24"
            />
          </label>
          <div className="flex gap-2 justify-end">
            <button
              type="button"
              className="btn-standard px-4"
              onClick={() => {
                props.closeHook();
              }}
            >
              OK
            </button>
            {isExternalSub && (
              <button
                type="button"
                onClick={async () => {
                  await API.saveSubDelay();
                  await refreshSubDelay();
                }}
                className="btn-standard px-4"
                title="Permanently save subtitle changes to file"
              >
                Save Subtitle Changes
              </button>
            )}
          </div>
        </div>
      </div>
    );
  }
  return null;
}

export { SettingsPopup };
