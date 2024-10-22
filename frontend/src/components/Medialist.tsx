import { API, type MediaItem } from "../lib/API";

function Medialist(props: {
  mediaItems: MediaItem[];
  subtitle: string;
}) {
  return (
    <div className="flex flex-col mb-2 pl-2 whitespace-pre=wrap">
      <h1 className="m-y-2/100 text-lg font-bold">{props.subtitle}</h1>
      {props.mediaItems.length !== 0
        ? props.mediaItems.map((mediaItem, i: number) => (
            <input
              type="button"
              className="text-left hover:bg-cyan-200 p-1 border-lightgray border-t-2 whitespace-normal break-words"
              key={props.subtitle + mediaItem.ID}
              onClick={(e) => API.setMedia(mediaItem)}
              value={mediaItem.Path}
            />
          ))
        : null}
    </div>
  );
}

export { Medialist };
