import { useCallback, useRef } from "react";

interface UseLongPressOptions {
  onShortClick: () => void;
  onLongPress: () => void;
  duration?: number;
}

function useLongPress({
  onShortClick,
  onLongPress,
  duration = 400,
}: UseLongPressOptions) {
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const isLongPressRef = useRef(false);
  const startX = useRef<number | null>(null);
  const startY = useRef<number | null>(null);
  const touchId = useRef<number | null>(null);
  const movedRef = useRef(false);
  const MOVE_THRESHOLD = 10; // pixels

  const handleMouseDown = useCallback(() => {
    isLongPressRef.current = false;

    timeoutRef.current = setTimeout(() => {
      isLongPressRef.current = true;
      onLongPress();
    }, duration);
  }, [duration, onLongPress]);

  const handleMouseUp = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    if (!isLongPressRef.current) {
      onShortClick();
    }

    isLongPressRef.current = false;
  }, [onShortClick]);

  const handleMouseLeave = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    isLongPressRef.current = false;
  }, []);

  const handleTouchStart = useCallback(
    (e: TouchEvent) => {
      if (!e.touches || e.touches.length === 0) return;
      // track the first touch only
      const t = e.touches[0];
      touchId.current = t.identifier;
      startX.current = t.clientX;
      startY.current = t.clientY;
      movedRef.current = false;
      isLongPressRef.current = false;

      // make touch slightly less sensitive: allow a longer duration
      const touchDuration = Math.max(duration, 600);

      timeoutRef.current = setTimeout(() => {
        isLongPressRef.current = true;
        onLongPress();
      }, touchDuration);
    },
    [duration, onLongPress],
  );

  const handleTouchMove = useCallback((e: TouchEvent) => {
    if (!touchId.current) return;
    for (let i = 0; i < e.touches.length; i++) {
      const t = e.touches[i];
      if (t.identifier === touchId.current) {
        const dx = Math.abs((startX.current ?? 0) - t.clientX);
        const dy = Math.abs((startY.current ?? 0) - t.clientY);
        if (dx > MOVE_THRESHOLD || dy > MOVE_THRESHOLD) {
          movedRef.current = true;
          if (timeoutRef.current) {
            clearTimeout(timeoutRef.current);
            timeoutRef.current = null;
          }
        }
        break;
      }
    }
  }, []);

  const handleTouchEnd = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    // if movement cancelled the long press, don't trigger short click
    if (!isLongPressRef.current && !movedRef.current) {
      onShortClick();
    }

    isLongPressRef.current = false;
    touchId.current = null;
    startX.current = null;
    startY.current = null;
    movedRef.current = false;
  }, [onShortClick]);

  const handleTouchCancel = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    isLongPressRef.current = false;
    touchId.current = null;
    startX.current = null;
    startY.current = null;
    movedRef.current = false;
  }, []);

  return {
    onMouseDown: handleMouseDown,
    onMouseUp: handleMouseUp,
    onMouseLeave: handleMouseLeave,
    onTouchStart: handleTouchStart,
    onTouchMove: handleTouchMove,
    onTouchEnd: handleTouchEnd,
    onTouchCancel: handleTouchCancel,
  };
}

export { useLongPress };
