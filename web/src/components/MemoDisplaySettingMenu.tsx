import { Option, Select } from "@mui/joy";
import clsx from "clsx";
import { Settings2Icon } from "lucide-react";
import { useMemoFilterStore } from "@/store/v1";
import { useTranslate } from "@/utils/i18n";
import { Popover, PopoverContent, PopoverTrigger } from "./ui/Popover";

interface Props {
  className?: string;
}

const MemoDisplaySettingMenu = ({ className }: Props) => {
  const t = useTranslate();
  const memoFilterStore = useMemoFilterStore();
  const isApplying = Boolean(memoFilterStore.orderByTimeAsc) !== false;

  return (
    <Popover>
      <PopoverTrigger
        className={clsx(className, isApplying ? "text-teal-600 bg-teal-50 dark:text-teal-500 dark:bg-teal-900 rounded-sm" : "opacity-40")}
      >
        <Settings2Icon className="w-4 h-auto shrink-0" />
      </PopoverTrigger>
      <PopoverContent align="end" alignOffset={-12} sideOffset={14}>
        <div className="flex flex-col gap-2">
          <div className="w-full flex flex-row justify-between items-center">
            <span className="text-sm shrink-0 mr-3">{t("memo.order-by")}</span>
            <Select value="displayTime">
              <Option value={"displayTime"}>{t("memo.display-time")}</Option>
            </Select>
          </div>
          <div className="w-full flex flex-row justify-between items-center">
            <span className="text-sm shrink-0 mr-3">{t("memo.direction")}</span>
            <Select value={memoFilterStore.orderByTimeAsc} onChange={(_, value) => memoFilterStore.setOrderByTimeAsc(Boolean(value))}>
              <Option value={false}>{t("memo.direction-desc")}</Option>
              <Option value={true}>{t("memo.direction-asc")}</Option>
            </Select>
          </div>
        </div>
      </PopoverContent>
    </Popover>
  );
};

export default MemoDisplaySettingMenu;
