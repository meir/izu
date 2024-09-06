local formatter = {}

-- Super + { a, b } + XF68Media{Play,Pause}
-- ________________________________________
-- In: ["Super", [ "a", "b" ], [ "XF68Media", [ "Play", "Pause" ]]]
-- Out: "Super + { a, b } + XF68Media{Play,Pause}
function formatter.base (parts)

end

-- Super + { a, b } + XF68Media{Play,Pause}
--         ^______^
-- In: [ "a", "b" ]
-- Out: [ "a", "b" ]
function formatter.multiple (parts)

end

-- Super + { a, b } + XF68Media{Play,Pause}
-- _____     _  _     _____________________
-- In: [ "Super" ], [ "a" ], [ "b" ], [ "XF68Media", [ "Play", "Pause" ]]
-- Out: [ "Super" ], [ "a" ], [ "b" ], [ "XF68Media", [ "Play", "Pause" ]]
function formatter.single (parts)

end

-- Super + { a, b } + XF68Media{Play,Pause}
--                             ^__________^
-- In: [ "Play", "Pause" ]
function formatter.single_part (parts)

end

-- Super + { a, b } + XF68Media{Play,Pause}
-- ^^^^^     ^  ^     ^^^^^^^^^ ^^^^ ^^^^^
-- In: "Super", "a", "b", "XF68Media", "Play", "Pause"
function formatter.string (parts)

end

return formatter
