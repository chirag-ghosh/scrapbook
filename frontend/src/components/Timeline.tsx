import { useEffect, useRef, useState } from "react";
import { IPhoto } from "../types/Photo";
import { BACKEND_URL, TIMELINE_IMAGE_WIDTH } from "../constants";

const getFormattedDate = (rawDate: Date) => {
  const date = new Date(rawDate);
  return (
    date.toLocaleDateString("en-GB") + " " + date.toLocaleTimeString("en-GB")
  );
};

const getPhotosInARow = () => {
    if (window.innerWidth < TIMELINE_IMAGE_WIDTH) {
        return 1;
    } else {
        return Math.floor(window.innerWidth / TIMELINE_IMAGE_WIDTH);
    }
}

const resizePhotos = (photos: IPhoto[], photosInARow: number, windowWidth: number) : IPhoto[] => {
    const rawResizedPhotos = photos.map((photo) => {
        return {
            ...photo,
            adjusted_width: TIMELINE_IMAGE_WIDTH,
            adjusted_height: (photo.height / photo.width) * TIMELINE_IMAGE_WIDTH,
        };
    });

    const resizedPhotos : IPhoto[] = [];
    for(let i = 0; i < rawResizedPhotos.length; i += photosInARow) {
        const row = rawResizedPhotos.slice(i, i + photosInARow);
        const rowAvgHeight = row.reduce((acc, photo) => acc + photo.adjusted_height, 0) / photosInARow;
        const rowWidth = row.reduce((acc, photo) => acc + photo.adjusted_width * (rowAvgHeight / photo.adjusted_height), 0);
        const rowHeight = rowAvgHeight * (windowWidth / rowWidth);
        row.forEach(photo => {
            resizedPhotos.push({
                ...photo,
                adjusted_height: rowHeight,
                adjusted_width: (photo.adjusted_width / photo.adjusted_height) * rowHeight,
            });
        });
    }

    return resizedPhotos;
}

function Timeline() {
  const [photosInARow, setPhotosInARow] = useState(getPhotosInARow());
  const [photos, setPhotos] = useState<IPhoto[]>([]);
  const [adjustedPhotos, setAdjustedPhotos] = useState<IPhoto[]>([]);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [hasNextPage, setHasNextPage] = useState(true);

  useEffect(() => {
    const handleResize = () => {
      setPhotosInARow(getPhotosInARow());
      setAdjustedPhotos(resizePhotos(photos, getPhotosInARow(), window.innerWidth))
    };
    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize);
    };
  }, [photos]);

  const observerTarget = useRef(null);

  const fetchPhotos = async () => {
    if (loading || !hasNextPage) return;
    setLoading(true);

    try {
      const response = await fetch(`${BACKEND_URL}/timeline?page=${page}&limit=${20}`);
      const data = await response.json();

      if (data == null || data.length === 0) {
        setHasNextPage(false);
      } else {
        setPhotos((prevPhotos) => [...prevPhotos, ...data]);
        setAdjustedPhotos((prevPhotos) => [...prevPhotos, ...resizePhotos(data, photosInARow, window.innerWidth)]);
        setPage((prevPage) => prevPage + 1);
      }
    } catch (error) {
      console.log(error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const observe = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.intersectionRatio > 0 && entry.intersectionRatio <= 1) {
            fetchPhotos();
          }
        });
      },
      {
        rootMargin: "0px 0px 200px 0px",
      }
    );

    if (observerTarget.current) {
      observe.observe(observerTarget.current);
    }

    return () => {
      if (observerTarget.current) {
        observe.unobserve(observerTarget.current);
      }
    };
  }, [observerTarget, photos, page, loading, hasNextPage]);

  return (
    <div className="timeline">
      {adjustedPhotos.map((photo) => (
        <div key={photo.id} className="timeline-item" style={{'width': photo.adjusted_width, 'height': photo.adjusted_height}}>
          <img
            src={`${BACKEND_URL}/photo/${photo.id}/serve`}
            alt={photo.name}
          />
          <div className="metadata">
            <p className="name">{photo.name}</p>
            <p className="date">{getFormattedDate(photo.captured_at)}</p>
          </div>
        </div>
      ))}
      {loading && hasNextPage && <p className="loading">Loading...</p>}
      <div ref={observerTarget}></div>
    </div>
  );
}

export default Timeline;
