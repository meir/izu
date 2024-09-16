local formatter = {}
local izu = izu

izu.registerKeycode({
  "Super",
  "Alt",
  "Shift",
})

function formatter.keybind (parts)
  local bind = parts[1]
  local command = parts[2]
  return bind .. "\n  " .. command
end

function formatter.command(parts)
  return table.concat(parts, "")
end

-- Super + { a, b } + XF68Media{Play,Pause}
-- ________________________________________
function formatter.base (parts)
  return table.concat(parts, " + ")
end

-- Super + { a, b } + XF68Media{Play,Pause}
--         ^______^
function formatter.multiple (parts)
  return "{ " .. table.concat(parts, ", ") .. " }"
end

-- Super + { a, b } + XF68Media{Play,Pause}
-- _____     _  _     _____________________
function formatter.single (parts)
  return table.concat(parts, "")
end

-- Super + { a, b } + XF68Media{Play,Pause}
--                             ^__________^
function formatter.single_part (parts)
  return "{" .. table.concat(parts, ",") .. "}"
end

-- Super + { a, b } + XF68Media{Play,Pause}
-- ^^^^^     ^  ^     ^^^^^^^^^ ^^^^ ^^^^^
function formatter.string (part, section)
  return part
end

return formatter
