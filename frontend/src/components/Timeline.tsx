import { useCallback, useEffect, useRef, useState } from "react";
import { IPhoto } from "../types/Photo";
import { BACKEND_URL } from "../constants";

const getFormattedDate = (rawDate: Date) => {
    const date = new Date(rawDate)
    return date.toLocaleDateString('en-GB') + ' ' + date.toLocaleTimeString('en-GB');
}

function Timeline() {
    const [photos, setPhotos] = useState<IPhoto[]>([]);
    const [page, setPage] = useState(1);
    const [loading, setLoading] = useState(false);
    const [hasNextPage, setHasNextPage] = useState(true);

    const observerTarget = useRef(null);

    const fetchPhotos = async () => {
        if(loading || !hasNextPage) return;
        setLoading(true);

        try {
            const response = await fetch(`${BACKEND_URL}/timeline?page=${page}`)
            const data = await response.json()

            if(data == null || data.length === 0) {
                setHasNextPage(false)
            } else {
                setPhotos(prevPhotos => [...prevPhotos, ...data])
                setPage(prevPage => prevPage + 1)
            }
        } catch(error) {
            console.log(error)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        const observe = new IntersectionObserver((entries) => {
            console.log(entries);
            entries.forEach(entry => {
                if(entry.intersectionRatio > 0 && entry.intersectionRatio <= 1) {
                    fetchPhotos()
                }
            })
        }, {
            rootMargin: "0px 0px 200px 0px"
        })

        if(observerTarget.current) {
            observe.observe(observerTarget.current);
        }

        return () => {
            if(observerTarget.current) {
                observe.unobserve(observerTarget.current);
            }
        }
    }, [observerTarget, photos, page, loading, hasNextPage]);

    return (
        <div className='timeline'>
            {photos.map((photo) => (
                <div key={photo.id} className='timeline-item'>
                    <img src={`${BACKEND_URL}/photo/${photo.id}/serve`} alt={photo.name} />
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