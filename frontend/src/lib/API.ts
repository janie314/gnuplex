interface MediaDir {
  Path: string;
  LastScanned: string;
}

interface MediaItem {
  ID: number;
  Path: string;
  LastPlayed: string;
}

interface MediaItemRes {
  res: MediaItem[];
  count: number;
}

interface FileExtension {
  ID: number;
  Extension: string;
}

interface SubTrack {
  id: number;
  title: string;
  selected: boolean;
}

interface Version {
  version: string;
  source_hash: string;
}

class API {
  public static async play() {
    return await fetch("/api/play", { method: "POST" });
  }

  public static async pause() {
    return await fetch("/api/pause", { method: "POST" });
  }

  public static async getPos(): Promise<number> {
    return await fetch("/api/pos")
      .then((res) => res.json())
      .catch((e) => {
        console.error("failed to get pos", e);
        return 0;
      });
  }

  public static async getTimeRemaining(): Promise<number> {
    return await fetch("/api/timeremaining")
      .then((res) => res.json())
      .catch((e) => {
        console.error("failed to get time remaining", e);
        return 0;
      });
  }

  public static async getVersion() {
    return (await fetch("/api/version").then((res) => res.json())) as Version;
  }

  public static async setPos(pos: number) {
    return await fetch("/api/pos", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ pos }),
    }).then((res) => res.json());
  }

  public static async getPaused(): Promise<boolean> {
    return await fetch("/api/paused")
      .then((res) => res.json())
      .catch((e) => {
        console.error("failed to get paused state", e);
        return 0;
      });
  }

  public static async getVol(): Promise<number> {
    return await fetch("/api/vol")
      .then((res) => res.json())
      .catch((e) => {
        console.error("failed to get vol", e);
        return 0;
      });
  }

  public static async setVol(vol: number) {
    return await fetch("/api/vol", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ vol }),
    }).then((res) => res.json());
  }

  public static async getNowPlaying() {
    return (await fetch("/api/nowplaying").then((res) =>
      res.json(),
    )) as MediaItem | null;
  }

  public static async setNowPlaying(mediaItem: MediaItem) {
    return await fetch("/api/nowplaying", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ id: mediaItem.ID }),
    });
  }

  public static async cast(url: string, temp: boolean) {
    return await fetch("/api/cast", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ url, temp }),
    });
  }

  public static async getMediadirs(): Promise<MediaDir[]> {
    return await fetch("/api/mediadirs")
      .then((res) => res.json())
      .catch((e) => {
        console.error(e);
        return [];
      });
  }

  public static async setMediadirs(mediadirs: string[]) {
    return await fetch("/api/mediadirs", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(mediadirs),
    });
  }

  public static async getFileExts() {
    return await fetch("/api/file_exts")
      .then((res) => res.json())
      .then((data: FileExtension[]) => {
        return data.sort((a, b) =>
          a.Extension.toLowerCase() < b.Extension.toLowerCase() ? -1 : 1,
        );
      });
  }

  public static async setFileExts(file_exts: string[]) {
    return await fetch("/api/file_exts", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(file_exts),
    });
  }

  public static async getMediaItems(search: string, paginationOffset: number) {
    const param = search || "";
    return (await fetch(
      `/api/mediaitems?search=${encodeURIComponent(param)}&offset=${encodeURIComponent(paginationOffset)}`,
    ).then((res) => res.json())) as MediaItemRes;
  }

  public static async scanLib() {
    return await fetch("/api/scanlib", { method: "POST" });
  }

  public static async getLast25Played() {
    return (await fetch("/api/last25").then((res) =>
      res.json(),
    )) as MediaItem[];
  }

  public static async getSubTracks() {
    return (await fetch("/api/sub").then((res) => res.json())) as
      | SubTrack[]
      | null;
  }

  public static async setSubTrack(id: number) {
    const body = id === -1 ? { visible: false } : { visible: true, id };
    return await fetch("/api/sub", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    });
  }
}

export { type MediaDir, type MediaItem, type SubTrack, API };
