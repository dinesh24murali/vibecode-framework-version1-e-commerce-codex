'use client'
import { useState } from 'react'
import { adminPresignImageUpload } from '@/lib/api/generated/admin-upload/admin-upload'
import { ImagePresignRequestContentType } from '@/lib/api/generated/model'
import { putToPresignedUrl } from '@/lib/api/client'
import { errorMessage } from '@/lib/auth/helpers'

export interface UploadedImage {
  url: string
  alt_text: string | null
}

const allowed = Object.values(ImagePresignRequestContentType) as string[]

/** Uploads images straight to S3 via a presigned PUT URL and reports the final public URLs. */
export function ProductImageUploader({
  images,
  onChange,
}: {
  images: UploadedImage[]
  onChange: (images: UploadedImage[]) => void
}) {
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function onFile(e: React.ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0]
    e.target.value = ''
    if (!file) return
    if (!allowed.includes(file.type)) {
      setError('Only JPEG, PNG or WebP images are allowed')
      return
    }
    setBusy(true)
    setError(null)
    try {
      const presign = await adminPresignImageUpload({
        filename: file.name,
        content_type: file.type as ImagePresignRequestContentType,
        file_size_bytes: file.size,
      })
      if (presign.status !== 200) {
        setError(errorMessage(presign, 'Could not start the upload'))
        return
      }
      const ok = await putToPresignedUrl(presign.data.presigned_url, file)
      if (!ok) {
        setError('Upload to storage failed')
        return
      }
      onChange([...images, { url: presign.data.public_url, alt_text: null }])
    } catch {
      setError('Upload failed. Please try again.')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="space-y-2">
      <span className="block text-sm font-medium">Images</span>
      <ul className="space-y-1 text-sm">
        {images.map((img, i) => (
          <li key={img.url} className="flex items-center justify-between gap-2">
            <span className="truncate">{img.url}</span>
            <button
              type="button"
              className="text-red-600 underline"
              onClick={() => onChange(images.filter((_, j) => j !== i))}
            >
              Remove
            </button>
          </li>
        ))}
      </ul>
      <input
        type="file"
        accept={allowed.join(',')}
        aria-label="Upload image"
        disabled={busy}
        onChange={(e) => void onFile(e)}
      />
      {error && (
        <p role="alert" className="text-sm text-red-600">
          {error}
        </p>
      )}
    </div>
  )
}
