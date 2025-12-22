import "./Popup.css";
import { API, type MediaItem } from "../lib/API";

function QueuePopup(props: {
  visible: boolean;
  mediaItem: MediaItem | null;
  closeHook: () => void;
}) {
  if (props.visible) {
    return (
      <div className="popup bg-white dark:bg-stone-800 m-5">
        <div className="flex flex-row">
          <input
            type="button"
            value="Play Next"
            className="btn-standard mr-1"
            onClick={() => {
              if (props.mediaItem) {
                API.playMedia(props.mediaItem, true, false);
              }
            }}
          />
          <input
            type="button"
            value="Play Last"
            className="btn-standard mr-1"
            onClick={() => {
              if (props.mediaItem) {
                API.playMedia(props.mediaItem, false, true);
              }
            }}
          />
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

export { QueuePopup };
