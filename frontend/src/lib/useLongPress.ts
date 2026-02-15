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
  const isTouchRef = useRef(false);

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
    isTouchRef.current = false;
  }, [onShortClick]);

  const handleMouseLeave = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    isLongPressRef.current = false;
  }, []);

  const handleTouchStart = useCallback(
    (e: React.TouchEvent) => {
      isTouchRef.current = true;
      isLongPressRef.current = false;

      timeoutRef.current = setTimeout(() => {
        isLongPressRef.current = true;
        onLongPress();
      }, duration);
    },
    [duration, onLongPress],
  );

  const handleTouchEnd = useCallback(
    (e: React.TouchEvent) => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }

      if (!isLongPressRef.current && isTouchRef.current) {
        e.preventDefault();
        onShortClick();
      }

      isLongPressRef.current = false;
      isTouchRef.current = false;
    },
    [onShortClick],
  );

  const handleTouchCancel = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    isLongPressRef.current = false;
    isTouchRef.current = false;
  }, []);

  return {
    onMouseDown: handleMouseDown,
    onMouseUp: handleMouseUp,
    onMouseLeave: handleMouseLeave,
    onTouchStart: handleTouchStart,
    onTouchEnd: handleTouchEnd,
    onTouchCancel: handleTouchCancel,
  };
}

export { useLongPress };
