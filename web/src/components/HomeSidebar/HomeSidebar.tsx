import clsx from "clsx";
import { useLocation } from "react-router-dom";
import useDebounce from "react-use/lib/useDebounce";
import SearchBar from "@/components/SearchBar";
import UserStatisticsView from "@/components/UserStatisticsView";
import useCurrentUser from "@/hooks/useCurrentUser";
import { useMemoList, useMemoMetadataStore } from "@/store/v1";
import TagsSection from "./TagsSection";
import useCurrentNest from "@/hooks/useCurrentNest";

interface Props {
  className?: string;
}

const HomeSidebar = (props: Props) => {
  const location = useLocation();
  const user = useCurrentUser();
  const memoList = useMemoList();
  const memoMetadataStore = useMemoMetadataStore();
  const nest = useCurrentNest();

  useDebounce(
    async () => {
      await memoMetadataStore.fetchMemoMetadata({ user, location, nest });
    },
    300,
    [memoList.size(), user, location.pathname, nest],
  );

  return (
    <aside
      className={clsx(
        "relative w-full h-auto max-h-screen overflow-auto hide-scrollbar flex flex-col justify-start items-start",
        props.className,
      )}
    >
      <SearchBar />
      <UserStatisticsView />
      <TagsSection />
    </aside>
  );
};

export default HomeSidebar;
