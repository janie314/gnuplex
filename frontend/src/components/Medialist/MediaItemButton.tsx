import { useLongPress } from "@uidotdev/usehooks";
import { API, type MediaItem } from "../../lib/API";

function MediaItemButton(props: {
  mediaItem: MediaItem;
  setQueueingTargetMediaItem: React.Dispatch<
    React.SetStateAction<MediaItem | null>
  >;
}) {
  const longPress = useLongPress(() => {}, {
    onFinish: () => {
      props.setQueueingTargetMediaItem(props.mediaItem);
    },
    onCancel: () => {
      API.playMedia(props.mediaItem, true, false);
    },
    threshold: 500,
  });

  return (
    <input
      type="button"
      className="inline-block w-screen max-w-screen-md text-left hover:bg-cyan-200 hover:dark:bg-cyan-700 p-1 border-lightgray dark:border-stone-800 border-t-2 whitespace-break-spaces break-words dark:bg-stone-950 dark:text-slate-200"
      value={props.mediaItem.Path}
      {...longPress}
    />
  );
}

export { MediaItemButton };
