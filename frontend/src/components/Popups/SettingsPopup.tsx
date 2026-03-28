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

  useEffect(() => {
    API.getSubDelay()
      .then(setSubDelay)
      .catch(() => setSubDelay(0));
  }, [props.visible]);

  if (props.visible) {
    return (
      <div className="popup bg-white dark:bg-stone-800 m-5 min-w-64 p-12">
        <div className="flex flex-col gap-3">
          <label className="flex items-center gap-2 text-black dark:text-white text-sm">
            Sub Delay (s):
            <input
              type="number"
              step="0.1"
              value={subDelay}
              onChange={(e) => setSubDelay(parseFloat(e.target.value) || 0)}
              onBlur={() => API.setSubDelay(subDelay)}
              className="btn-standard w-24"
            />
          </label>
          <div className="flex items-center gap-2">
            <span className="text-black dark:text-white text-sm whitespace-nowrap">
              Seek Subtitle:
            </span>
            <div className="flex gap-2">
              <button
                onClick={() => API.subSeek(-1)}
                className="btn-standard px-3"
                title="Previous subtitle"
              >
                <svg
                  viewBox="0 0 24 24"
                  className="w-5 h-5"
                  fill="currentColor"
                >
                  <path d="M6 6h2v12H6V6zm3.5 6l8.5 6V6l-8.5 6z" />
                </svg>
              </button>
              <button
                onClick={() => API.subSeek(0)}
                className="btn-standard px-3"
                title="Jump to current subtitle start"
              >
                <svg
                  viewBox="0 0 24 24"
                  className="w-5 h-5"
                  fill="currentColor"
                >
                  <path d="M12 5V1L7 6l5 5V7c3.31 0 6 2.69 6 6s-2.69 6-6 6-6-2.69-6-6H4c0 4.42 3.58 8 8 8s8-3.58 8-8-3.58-8-8-8z" />
                </svg>
              </button>
              <button
                onClick={() => API.subSeek(1)}
                className="btn-standard px-3"
                title="Next subtitle"
              >
                <svg
                  viewBox="0 0 24 24"
                  className="w-5 h-5"
                  fill="currentColor"
                >
                  <path d="M6 18l8.5-6L6 6v12zm2-12v12l6.5-6L8 6zm8 0v12h2V6h-2z" />
                </svg>
              </button>
            </div>
          </div>
          <select
            className="btn-standard w-full"
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
          <button
            className="btn-standard self-start px-4"
            onClick={() => {
              props.closeHook();
            }}
          >
            OK
          </button>
        </div>
      </div>
    );
  }
  return null;
}

export { SettingsPopup };
