export interface IPhoto {
    id: number
    name: string
    file_dir: string
    camera_make: string
    camera_model: string
    lens_id: string
    width: number
    adjusted_width: number
    height: number
    adjusted_height: number
    focal_length: number
    aperture: number
    shutter_speed: string
    iso: number
    captured_at: Date
}
