import { API, type MediaItem } from "../lib/API";
import { PageSelector } from "./Medialist/PageSelector";
import queue from "../assets/queue.svg";

function Medialist(props: {
  mediaItems: (MediaItem | null)[];
  subtitle: string;
  mediaItemCount: number | null;
  paginationOffset: number | null;
  setPaginationOffset: React.Dispatch<React.SetStateAction<number>> | null;
}) {
  if (props.mediaItems.length === 0) {
    return null;
  }
  return (
    <div className="w-full flex flex-col mb-2 pl-2 whitespace-pre=wrap">
      <div className="flex flex-row align-center mb-1">
        <h1 className="header">{props.subtitle}</h1>
        <PageSelector
          mediaItemCount={props.mediaItemCount}
          paginationOffset={props.paginationOffset}
          setPaginationOffset={props.setPaginationOffset}
        />
      </div>
      {props.mediaItems
        .filter((mediaItem) => mediaItem !== null)
        .map((mediaItem, i: number) => (
          <div
            key={props.subtitle + mediaItem.ID}
            className="flex flex-row items-center"
          >
            <input
              type="button"
              className="inline-block w-screen max-w-screen-md text-left hover:bg-cyan-200 hover:dark:bg-cyan-700 p-1 border-lightgray dark:border-stone-800 border-t-1 border-b-1 whitespace-break-spaces break-words dark:bg-stone-950 dark:text-slate-200"
              onClick={(e) => API.setNowPlaying(mediaItem)}
              value={mediaItem.Path}
            />
            <button
              type="button"
              className="p-3 inline-flex items-center text-xl text-bold h-full hover:dark:bg-cyan-700 dark:text-white"
              onClick={() => API.pause()}
            >
              ...
            </button>
          </div>
        ))}
    </div>
  );
}

export { Medialist };
