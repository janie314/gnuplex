import "./Popup.css";
import { API, type MediaItem } from "../lib/API";

function QueuePopup(props: {
  visible: boolean;
  setQueuePopup: React.Dispatch<React.SetStateAction<boolean>>;
  mediaItem: MediaItem | null;
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
                API.setNowPlaying(props.mediaItem);
              }
            }}
          />{" "}
          <input
            type="button"
            value="Queue Next"
            className="btn-standard mr-1"
            onClick={() => {}}
          />
          <input
            type="button"
            value="Queue Last"
            className="btn-standard mr-1"
            onClick={() => {}}
          />
          <input
            type="button"
            value="Cancel"
            className="btn-standard"
            onClick={() => {
              setVol(props.currentVol);
              props.setVolInputPopup(false);
            }}
          />
        </div>
      </div>
    );
  }
  return null;
}

export { QueuePopup };
