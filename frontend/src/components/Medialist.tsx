import { API, type MediaItem } from "../lib/API";
import { PageSelector } from "./Medialist/PageSelector";

function Medialist(props: {
  mediaItems: (MediaItem | null)[];
  subtitle: string;
  mediaItemCount: number | null;
  paginationOffset: number | null;
  setPaginationOffset: React.Dispatch<React.SetStateAction<number>> | null;
}) {
  if (props.mediaItems === null || props.mediaItems.length === 0) {
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
          <input
            type="button"
            className="inline-block w-screen max-w-screen-md text-left hover:bg-cyan-200 hover:dark:bg-cyan-700 p-1 border-lightgray dark:border-stone-800 border-t-2 whitespace-break-spaces break-words dark:bg-stone-950 dark:text-slate-200"
            key={props.subtitle + mediaItem.ID}
            onClick={(e) => API.setNowPlaying(mediaItem)}
            value={mediaItem.Path}
          />
        ))}
    </div>
  );
}

export { Medialist };
