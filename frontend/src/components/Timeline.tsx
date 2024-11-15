import { useEffect, useState } from "react";
import { IPhoto } from "../types/Photo";
import { BACKEND_URL } from "../constants";

const getFormattedDate = (rawDate: Date) => {
    const date = new Date(rawDate)
    return date.toLocaleDateString('en-GB') + ' ' + date.toLocaleTimeString('en-GB');
}

function Timeline() {
    const [photos, setPhotos] = useState<IPhoto[]>([]);

    useEffect(() => {
        fetch(`${BACKEND_URL}/timeline`)
            .then((res) => res.json())
            .then((data) => {
                setPhotos(data);
            })
    }, []);

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
        </div>
    );
}

export default Timeline;