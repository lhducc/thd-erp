import { useState, useEffect, useRef } from "react";
import {
  MapContainer,
  TileLayer,
  Marker,
  Popup,
  useMap,
  useMapEvents,
} from "react-leaflet";
import { Search, Loader2, X } from "lucide-react";
import "leaflet/dist/leaflet.css";
import L from "leaflet";

// Import các component từ shadcn/ui
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { useDebounce } from "@uidotdev/usehooks";

// Sửa biểu tượng marker mặc định của Leaflet
delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png",
  iconUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png",
  shadowUrl:
    "https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-shadow.png",
});

// Component để cập nhật view của bản đồ
function MapUpdater({ center }) {
  const map = useMap();

  useEffect(() => {
    if (center) {
      map.flyTo(center, 15);
    }
  }, [center, map]);

  return null;
}

// Component để xử lý click trên bản đồ
function MapClickHandler({ onMapClick }) {
  useMapEvents({
    click: (e) => {
      onMapClick([e.latlng.lat, e.latlng.lng]);
    },
  });

  return null;
}

// Component chính
export default function Map({
  setCoordinates,
  coordinates,
}: {
  setCoordinates: Function;
  coordinates: { latitude: number; longitude: number };
}) {
  const [searchQuery, setSearchQuery] = useState("");
  const [suggestions, setSuggestions] = useState([]);
  const [selectedLocation, setSelectedLocation] = useState(null);
  const [markerPosition, setMarkerPosition] = useState(null);
  // const [coordinates, setCoordinates] = useState(null); // State để lưu tọa độ
  const [isLoading, setIsLoading] = useState(false);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const markerRef = useRef(null);
  const debouncedSearchTerm = useDebounce(searchQuery, 500);

  // Vị trí mặc định (Hà Nội)
  const defaultPosition = [21.0285, 105.8542];

  useEffect(() => {
    if (coordinates) {
      setMarkerPosition([coordinates.latitude, coordinates.longitude]);
    }
  }, []);

  // Cập nhật tọa độ khi markerPosition thay đổi
  useEffect(() => {
    if (markerPosition) {
      setCoordinates({
        lat: markerPosition[0],
        lng: markerPosition[1],
      });
    }
  }, [markerPosition]);

  // Tìm kiếm địa điểm với Nominatim API
  const searchLocations = async (query) => {
    if (query.length < 3) {
      setSuggestions([]);
      return;
    }

    setIsLoading(true);
    try {
      const response = await fetch(
        `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(
          query
        )}&format=json&limit=5`
      );
      const data = await response.json();
      setSuggestions(data);
      setShowSuggestions(true);
    } catch (error) {
      console.error("Lỗi khi tìm kiếm địa điểm:", error);
      setSuggestions([]);
    } finally {
      setIsLoading(false);
    }
  };

  // Xử lý thay đổi input tìm kiếm
  const handleSearchChange = (e) => {
    const query = e.target.value;
    setSearchQuery(query);

    // Debounce tìm kiếm
    // const timeoutId = setTimeout(() => {
    //   searchLocations(query);
    // }, 500);

    // return () => clearTimeout(timeoutId);
  };

  useEffect(() => {
    searchLocations(debouncedSearchTerm);
  }, [debouncedSearchTerm]);

  // Xử lý khi chọn một địa điểm
  const handleSelectLocation = (location) => {
    setSelectedLocation(location);
    const newPosition = [parseFloat(location.lat), parseFloat(location.lon)];
    setMarkerPosition(newPosition);
    setShowSuggestions(false);
    setSearchQuery(`${location.display_name}`);
  };

  // Xử lý khi marker được kéo
  const eventHandlers = {
    dragend() {
      const marker = markerRef.current;
      if (marker) {
        const newPosition = marker.getLatLng();
        setMarkerPosition([newPosition.lat, newPosition.lng]);
      }
    },
  };

  // Xử lý khi click vào bản đồ
  const handleMapClick = (position) => {
    setMarkerPosition(position);
  };

  // Xóa marker và reset tìm kiếm
  const handleClearSearch = () => {
    setSearchQuery("");
    setSelectedLocation(null);
    setMarkerPosition(null);
    setCoordinates(null);
    setSuggestions([]);
    setShowSuggestions(false);
  };

  return (
    <>
      <div className="relative">
        <div className="flex gap-2">
          <div className="relative flex-grow">
            <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
            <Input
              type="text"
              value={searchQuery}
              onChange={handleSearchChange}
              placeholder="Tìm kiếm địa điểm..."
              className="pl-8 pr-8"
            />
            {searchQuery && (
              <Button
                variant="ghost"
                size="icon"
                className="absolute right-1 top-1 h-7 w-7"
                onClick={handleClearSearch}
              >
                <X className="h-4 w-4" />
              </Button>
            )}
          </div>
        </div>

        {/* Hiển thị gợi ý */}
        {showSuggestions && suggestions.length > 0 && (
          <Card className="absolute w-full mt-1 max-h-64 overflow-y-auto z-[9999]">
            <CardContent className="p-0">
              {suggestions.map((place) => (
                <div
                  key={place.place_id}
                  className="p-3 hover:bg-accent cursor-pointer border-b"
                  onClick={() => handleSelectLocation(place)}
                >
                  <div className="font-medium">{place.display_name}</div>
                  <div className="text-xs text-muted-foreground mt-1 flex items-center gap-1">
                    <Badge variant="outline" className="text-xs">
                      {place.type}
                    </Badge>
                    {place.type !== place.class && (
                      <Badge variant="outline" className="text-xs">
                        {place.class}
                      </Badge>
                    )}
                  </div>
                </div>
              ))}
            </CardContent>
          </Card>
        )}

        {isLoading && (
          <div className="absolute right-12 top-2.5">
            <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
          </div>
        )}
      </div>

      {/* Hiển thị thông tin vị trí đã chọn */}
      {/* {coordinates && (
        <Alert className="bg-primary/10 border-primary/20">
          <div className="flex flex-col gap-2">
            <div className="flex items-center gap-2">
              <MapPin className="h-4 w-4 text-primary" />
              <span className="font-medium">Vị trí đã chọn</span>
            </div>
            <Separator />
            <div className="grid grid-cols-2 gap-2 mt-1">
              <div>
                <Label className="text-xs text-muted-foreground">
                  Vĩ độ (Latitude)
                </Label>
                <div className="font-mono text-sm">
                  {coordinates.lat.toFixed(6)}
                </div>
              </div>
              <div>
                <Label className="text-xs text-muted-foreground">
                  Kinh độ (Longitude)
                </Label>
                <div className="font-mono text-sm">
                  {coordinates.lng.toFixed(6)}
                </div>
              </div>
            </div>
            <AlertDescription className="text-xs text-muted-foreground mt-1">
              Bạn có thể click vào bản đồ hoặc kéo marker để điều chỉnh vị trí
            </AlertDescription>
          </div>
        </Alert>
      )} */}

      {/* Bản đồ */}
      <div className="h-96 rounded-lg overflow-hidden border">
        <MapContainer
          center={markerPosition || defaultPosition}
          zoom={13}
          style={{ height: "100%", width: "100%" }}
        >
          <TileLayer
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          />

          {markerPosition && (
            <Marker
              position={markerPosition}
              draggable={true}
              eventHandlers={eventHandlers}
              ref={markerRef}
            >
              <Popup>
                <div>
                  <div className="font-medium">Vị trí đã chọn</div>
                  <div className="text-xs">
                    {markerPosition[0].toFixed(6)},{" "}
                    {markerPosition[1].toFixed(6)}
                  </div>
                </div>
              </Popup>
            </Marker>
          )}

          <MapUpdater center={markerPosition} />
          <MapClickHandler onMapClick={handleMapClick} />
        </MapContainer>
      </div>
    </>
  );
}
