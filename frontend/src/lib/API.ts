interface PosResponse {
  pos: number;
  max_pos: number;
}

interface MediaStateResponse {
  pos: number;
  max_pos: number;
  vol: number;
  paused: boolean;
  media: string;
}

const common_options = {
  headers: {
    "Content-Type": "application/json",
  },
};

class APICall {
  /**
   * @returns GNUPlex version string.
   */
  public static async version(): Promise<string> {
    return await fetch(
      "/api/version",
      { ...common_options },
    ).then((res) => res.json());
  }

  /**
   * @returns Whether the video player is paused.
   */
  public static async paused(): Promise<boolean | null> {
    return await fetch(
      "/api/paused",
      { ...common_options },
    ).then((res) => res.json());
  }

  /**
   * @returns Full current state of the media player.
   */
  public static async mediastate(): Promise<MediaStateResponse | null> {
    return await fetch(
      "/api/mediastate",
      { ...common_options },
    ).then((res) => res.json());
  }

  /**
   * Toggles the video's play/pause status.
   *
   * @returns Whether the video player is paused following the toggle operation.
   */
  public static async toggle(): Promise<boolean> {
    return await fetch(
      "/api/toggle",
      { method: "POST", ...common_options },
    ).then((res) => res.json());
  }

  /**
   * @returns The video player's position (seconds).
   */
  public static async pos(): Promise<PosResponse | null> {
    return await fetch(
      "/api/pos",
      { ...common_options },
    ).then((res) => res.json());
  }

  /**
   * Sets the video's position, either with an absolute position or an increment.
   *
   * @param pos Position to set, either an increment or an absolute position.
   * @param inc Whether `pos` is an increment or absolute position.
   * @returns The video player's position (seconds) following the toggle inc operation.
   */
  public static async setPos(
    pos: number,
    inc: boolean,
  ) {
    return await fetch(
      "/api/pos",
      { method: "POST", ...common_options, body: JSON.stringify({ inc, pos }) },
    );
  }

  /**
   * @returns The video's volume (percentage, 0 - 100+).
   */
  public static async vol(): Promise<number | null> {
    return await fetch(
      "/api/vol",
      { ...common_options },
    ).then((res) => res.json());
  }

  /**
   * @param vol The video's volume (percentage, 0 - 100+).
   */
  public static async setVol(vol: number) {
    return await fetch(
      "/api/vol",
      { method: "POST", ...common_options, body: JSON.stringify({ vol }) },
    );
  }

  public static async media(): Promise<string | null> {
    return await fetch(
      "/api/media",
      { ...common_options },
    ).then((res) => res.json());
  }

  public static async setMedia(mediafile: string) {
    return await fetch(
      "/api/media",
      {
        method: "POST",
        ...common_options,
        body: JSON.stringify({ media: mediafile }),
      },
    );
  }

  public static async mediadirs(): Promise<string[]> {
    const res: string[] = await fetch(
      "/api/mediadirs",
    ).then((res) => res.json());
    return res.sort((a, b) => a.toLowerCase() < b.toLowerCase() ? -1 : 1);
  }

  public static async setMediadirs(mediadirs: string[]) {
    return await fetch(
      "/api/mediadirs",
      { method: "POST", ...common_options, body: JSON.stringify(mediadirs) },
    );
  }

  public static async fileExts(): Promise<string[]> {
    const res: string[] = await fetch(
      "/api/file_exts",
    ).then((res) => res.json());
    return res.sort((a, b) => a.toLowerCase() < b.toLowerCase() ? -1 : 1);
  }

  public static async setFileExts(file_exts: string[]) {
    return await fetch(
      "/api/file_exts",
      { method: "POST", ...common_options, body: JSON.stringify(file_exts) },
    );
  }

  public static async mediafiles() {
    return await fetch(`/api/medialist`).then((res) => res.json());
  }

  public static async setMediafiles() {
    return await fetch(`/api/medialist`, { method: "POST" });
  }

  public static async last25() {
    return await fetch(`/api/last25`).then((res) => res.json());
  }

  public static sleep(ms: number) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }
}

export { type PosResponse };
export { APICall };
