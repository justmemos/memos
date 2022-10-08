import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useAppSelector } from "../store";
import { locationService, memoService, userService } from "../services";
import useToggle from "../hooks/useToggle";
import closeSidebar from "../hooks/closeSidebar"
import Icon from "./Icon";
import "../less/tag-list.less";

interface Tag {
  key: string;
  text: string;
  subTags: Tag[];
}

const TagList = () => {
  const { t } = useTranslation();
  const { memos, tags: tagsText } = useAppSelector((state) => state.memo);
  const query = useAppSelector((state) => state.location.query);
  const [tags, setTags] = useState<Tag[]>([]);

  useEffect(() => {
    if (memos.length > 0) {
      memoService.updateTagsState();
    }
  }, [memos]);

  useEffect(() => {
    const sortedTags = Array.from(tagsText).sort();
    const root: KVObject<any> = {
      subTags: [],
    };

    for (const tag of sortedTags) {
      const subtags = tag.split("/");
      let tempObj = root;
      let tagText = "";

      for (let i = 0; i < subtags.length; i++) {
        const key = subtags[i];
        if (i === 0) {
          tagText += key;
        } else {
          tagText += "/" + key;
        }

        let obj = null;

        for (const t of tempObj.subTags) {
          if (t.text === tagText) {
            obj = t;
            break;
          }
        }

        if (!obj) {
          obj = {
            key,
            text: tagText,
            subTags: [],
          };
          tempObj.subTags.push(obj);
        }

        tempObj = obj;
      }
    }

    setTags(root.subTags as Tag[]);
  }, [tagsText]);

  return (
    <div className="tags-wrapper" onClick={() => closeSidebar()}>
      <p className="title-text">{t("common.tags")}</p>
      <div className="tags-container">
        {tags.map((t, idx) => (
          <TagItemContainer key={t.text + "-" + idx} tag={t} tagQuery={query?.tag} />
        ))}
        {!userService.isVisitorMode() && tags.length === 0 && <p className="tip-text">{t("tag-list.tip-text")}</p>}
      </div>
    </div>
  );
};

interface TagItemContainerProps {
  tag: Tag;
  tagQuery?: string;
}

const TagItemContainer: React.FC<TagItemContainerProps> = (props: TagItemContainerProps) => {
  const { tag, tagQuery } = props;
  const isActive = tagQuery === tag.text;
  const hasSubTags = tag.subTags.length > 0;
  const [showSubTags, toggleSubTags] = useToggle(false);

  const handleTagClick = () => {
    if (isActive) {
      locationService.setTagQuery(undefined);
    } else {
      locationService.setTagQuery(tag.text);
    }
  };

  const handleToggleBtnClick = (event: React.MouseEvent) => {
    event.stopPropagation();
    toggleSubTags();
  };

  return (
    <>
      <div className={`tag-item-container ${isActive ? "active" : ""}`} onClick={handleTagClick}>
        <div className="tag-text-container">
          <span className="icon-text">#</span>
          <span className="tag-text">{tag.key}</span>
        </div>
        <div className="btns-container">
          {hasSubTags ? (
            <span className={`action-btn toggle-btn ${showSubTags ? "shown" : ""}`} onClick={handleToggleBtnClick}>
              <Icon.ChevronRight className="icon-img" />
            </span>
          ) : null}
        </div>
      </div>

      {hasSubTags ? (
        <div className={`subtags-container ${showSubTags ? "" : "!hidden"}`}>
          {tag.subTags.map((st, idx) => (
            <TagItemContainer key={st.text + "-" + idx} tag={st} tagQuery={tagQuery} />
          ))}
        </div>
      ) : null}
    </>
  );
};

export default TagList;
