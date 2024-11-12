import { API, type MediaItem } from "../lib/API";

function Medialist(props: {
  mediaItems: (MediaItem | null)[];
  subtitle: string;
}) {
  if (props.mediaItems.length === 0) {
    return null;
  }
  return (
    <div className="flex flex-col mb-2 pl-2 whitespace-pre=wrap">
      <h1 className="m-y-2/100 text-lg font-bold">{props.subtitle}</h1>
      {props.mediaItems
        .filter((mediaItem) => mediaItem !== null)
        .map((mediaItem, i: number) => (
          <input
            type="button"
            className="text-left hover:bg-cyan-200 p-1 border-lightgray border-t-2 whitespace-normal break-words"
            key={props.subtitle + mediaItem.ID}
            onClick={(e) => API.setNowPlaying(mediaItem)}
            value={mediaItem.Path}
          />
        ))}
    </div>
  );
}

export { Medialist };
