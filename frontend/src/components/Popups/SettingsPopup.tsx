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
  if (props.visible) {
    return (
      <div className="popup bg-white dark:bg-stone-800 m-5">
        <div className="flex flex-col gap-2">
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
          <input
            type="button"
            value="Cancel"
            className="btn-standard"
            onClick={() => {
              props.closeHook();
            }}
          />
        </div>
      </div>
    );
  }
  return null;
}

export { SettingsPopup };
