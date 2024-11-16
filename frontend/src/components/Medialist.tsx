import { API, type MediaItem } from "../lib/API";
import { PageSelector } from "./Medialist/PageSelector";

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
        <h1 className="m-y-2/100 text-lg font-bold">{props.subtitle}</h1>
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
            className="inline-block w-screen max-w-screen-md text-left hover:bg-cyan-200 p-1 border-lightgray border-t-2 whitespace-break-spaces break-words"
            key={props.subtitle + mediaItem.ID}
            onClick={(e) => API.setNowPlaying(mediaItem)}
            value={mediaItem.Path}
          />
        ))}
    </div>
  );
}

export { Medialist };
