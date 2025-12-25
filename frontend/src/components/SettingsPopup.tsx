import "./Popup.css";
import { API } from "../lib/API";

function SettingsPopup(props: { visible: boolean; closeHook: () => void }) {
  if (props.visible) {
    return (
      <div className="popup bg-white dark:bg-stone-800 m-5">
        <div className="flex flex-col">
          <input
            type="button"
            value="Black and White Filter"
            className="btn-standard m-1 min-w-[11ch]"
            onClick={() => {
              API.setFilter("bw");
            }}
          />{" "}
          <input
            type="button"
            value="Grainy Filter"
            className="btn-standard m-1 min-w-[11ch]"
            onClick={() => {
              API.setFilter("grainy");
            }}
          />
          <input
            type="button"
            value="Reset Filter"
            className="btn-standard m-1 min-w-[11ch]"
            onClick={() => {
              API.setFilter("neutral");
            }}
          />
          <input
            type="button"
            value="Cancel"
            className="btn-standard m-1 min-w-[11ch]"
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
