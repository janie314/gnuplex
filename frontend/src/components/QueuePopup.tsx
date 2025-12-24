import "./Popup.css";
import { API, type MediaItem } from "../lib/API";

function QueuePopup(props: {
  visible: boolean;
  mediaItem: MediaItem | null;
  setQueueingTargetMediaItem: React.Dispatch<
    React.SetStateAction<MediaItem | null>
  >;
  closeHook: () => void;
}) {
  if (props.visible) {
    return (
      <div className="popup bg-white dark:bg-stone-800 m-5">
        <div className="flex flex-col">
          <input
            type="button"
            value="Play"
            className="btn-standard m-1 min-w-[11ch]"
            onClick={() => {
              if (props.mediaItem) {
                API.playMedia(props.mediaItem, true, false);
              }
              props.setQueueingTargetMediaItem(null);
            }}
          />
          <input
            type="button"
            value="Queue Next"
            className="btn-standard m-1 min-w-[11ch]"
            onClick={() => {
              if (props.mediaItem) {
                API.playMedia(props.mediaItem, true, false);
              }
              props.setQueueingTargetMediaItem(null);
            }}
          />
          <input
            type="button"
            value="Queue Last"
            className="btn-standard m-1 min-w-[11ch]"
            onClick={() => {
              if (props.mediaItem) {
                API.playMedia(props.mediaItem, false, true);
              }
              props.setQueueingTargetMediaItem(null);
            }}
          />
          <input
            type="button"
            value="Cancel"
            className="btn-standard m-1 min-w-[11ch]"
            onClick={() => {
              props.closeHook();
              props.setQueueingTargetMediaItem(null);
            }}
          />
        </div>
      </div>
    );
  }
  return null;
}

export { QueuePopup };
