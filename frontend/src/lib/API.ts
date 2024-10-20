interface IMPVRes {
  data?: number | string;
  request_id: number;
  error: string;
}

interface MediaDir {
  Path: string;
  LastScanned: string;
}

interface MediaItem {
  ID: number;
  Path: string;
  LastPlayed: string;
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

  public static async getTimeRemaining() {
    return await fetch("/api/timeremaining")
      .then((res) => res.json())
      .then((res: IMPVRes) => {
        if (res.data !== undefined) {
          // @ts-ignore
          return Math.floor(res.data);
        }
        return 0;
      });
  }

  public static async getVersion() {
    return (await fetch("/api/version").then((res) => res.json())) as string;
  }

  public static async setPos(pos: number) {
    return await fetch(`/api/pos?pos=${pos}`, { method: "POST" }).then((res) =>
      res.json(),
    );
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
    return await fetch(`/api/vol?vol=${vol}`, { method: "POST" }).then((res) =>
      res.json(),
    );
  }

  public static async getMedia() {
    return (await fetch("/api/media").then((res) => res.json())) as MediaItem;
  }

  public static async setMedia(mediaItem: MediaItem) {
    console.log("m", mediaItem);
    return await fetch("/api/media", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ id: mediaItem.ID }),
    });
  }

  public static async castMedia(url: string) {
    return await fetch("/api/cast", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ arg: url }),
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
      .then((data: string[]) => {
        return data.sort((a, b) =>
          a.toLowerCase() < b.toLowerCase() ? -1 : 1,
        );
      });
  }

  public static async setFileExts(file_exts: string[]) {
    return await fetch(
      `/api/file_exts?file_exts=${encodeURI(JSON.stringify(file_exts))}`,
      { method: "POST" },
    );
  }

  public static async getMediaItems() {
    return (await fetch("/api/mediaitems").then((res) =>
      res.json(),
    )) as MediaItem[];
  }

  public static async scanLib() {
    return await fetch("/api/scanlib", { method: "POST" });
  }

  public static async getLast25() {
    return (await fetch("/api/last25").then((res) =>
      res.json(),
    )) as MediaItem[];
  }

  public static sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }
}

export { type MediaDir, type MediaItem, API };
