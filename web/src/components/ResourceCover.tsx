import React from "react";
import Icon from "./Icon";

interface ResourceCoverProps {
  resource: Resource;
}

const ResourceCover = ({ resource }: ResourceCoverProps) => {
  switch (resource.type) {
    case "image/*":
      return <Icon.FileImage className="w-32 h-32 ml-auto mr-auto" />;
    case "video/*":
      return <Icon.FileVideo2 className="w-32 h-32 ml-auto mr-auto" />;
    case "audio/*":
      return <Icon.FileAudio className="w-32 h-32 ml-auto mr-auto" />;
    case "text/*":
      return <Icon.FileText className="w-32 h-32 ml-auto mr-auto" />;
    case "application/epub+zip":
      return <Icon.Book className="w-32 h-32 ml-auto mr-auto" />;
    case "application/pdf":
      return <Icon.Book className="w-32 h-32 ml-auto mr-auto" />;
    case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
      return <Icon.FileEdit className="w-32 h-32 ml-auto mr-auto" />;
    case "application/msword":
      return <Icon.FileEdit className="w-32 h-32 ml-auto mr-auto" />;
    default:
      return <Icon.File className="w-32 h-32 ml-auto mr-auto" />;
  }
};

export default React.memo(ResourceCover);
